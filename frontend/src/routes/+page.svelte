<script lang="ts">
	let url = '';
	let maxDepth = 3;
	let delay = 1000;
	let followExternal = false;
	let isLoading = false;
	let results: any = null;
	let error = '';

	const API_BASE = 'http://localhost:8080/api/v1';

	async function handleSubmit() {
		if (!url) {
			error = 'Please enter a URL';
			return;
		}

		isLoading = true;
		error = '';
		results = null;

		console.log('🚀 Starting website scraping...', { url, maxDepth, delay, followExternal });

		try {
			const response = await fetch(`${API_BASE}/scrape`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					url,
					maxDepth,
					delay,
					followExternal
				})
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Failed to scrape website');
			}

			if (data.success) {
				results = data;
				console.log('✅ Scraping completed!', data.stats);
			} else {
				throw new Error(data.error || 'Scraping failed');
			}
		} catch (err) {
			console.error('❌ Scraping error:', err);
			error = err instanceof Error ? err.message : 'An unknown error occurred';
		} finally {
			isLoading = false;
		}
	}

	function downloadAsJson() {
		if (!results?.pages) return;

		const filename = generateFileName(url, 'json');
		const blob = new Blob([JSON.stringify(results.pages, null, 2)], {
			type: 'application/json'
		});
		const downloadUrl = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = downloadUrl;
		a.download = filename;
		a.click();
		URL.revokeObjectURL(downloadUrl);
		console.log('💾 Downloaded JSON file:', filename);
	}

	function downloadAsMarkdown() {
		if (!results?.pages) return;

		const filename = generateFileName(url, 'md');
		let content = `# Website Content\n\n*Scraped on ${new Date().toLocaleString()}*\n\n---\n\n`;

		results.pages.forEach((page: any, index: number) => {
			if (page.error) {
				content += `## ❌ Error: ${page.url}\n\n**Error:** ${page.error}\n\n`;
				return;
			}

			content += `## 📄 Page ${index + 1}: ${page.title}\n\n`;
			content += `**URL:** ${page.url}  \n`;
			content += `**Depth:** ${page.depth}\n\n`;
			content += '---\n\n';
			content += page.markdown;
			content += '\n\n';

			if (index < results.pages.length - 1) {
				content += '---\n\n';
			}
		});

		const blob = new Blob([content], { type: 'text/markdown' });
		const downloadUrl = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = downloadUrl;
		a.download = filename;
		a.click();
		URL.revokeObjectURL(downloadUrl);
		console.log('💾 Downloaded Markdown file:', filename);
	}

	function generateFileName(websiteUrl: string, extension: string): string {
		try {
			const parsedUrl = new URL(websiteUrl);
			let siteName = parsedUrl.hostname.replace(/^www\./, ''); // Remove www. prefix

			// Clean up the site name to be filesystem-safe
			siteName = siteName.replace(/[^a-zA-Z0-9.-]/g, '-');

			// Generate timestamp
			const now = new Date();
			const timestamp = now
				.toISOString()
				.replace(/[:]/g, '-') // Replace colons with dashes
				.replace(/\..+/, '') // Remove milliseconds
				.replace('T', '_'); // Replace T with underscore

			return `${siteName}_${timestamp}.${extension}`;
		} catch (error) {
			// Fallback if URL parsing fails
			const timestamp = new Date()
				.toISOString()
				.replace(/[:]/g, '-')
				.replace(/\..+/, '')
				.replace('T', '_');

			return `website_${timestamp}.${extension}`;
		}
	}
</script>

<div class="min-h-screen bg-gray-50 px-4 py-8">
	<div class="mx-auto max-w-4xl">
		<!-- Header -->
		<div class="mb-8 text-center">
			<h1 class="mb-2 text-4xl font-bold text-gray-900">🔄 Website to Markdown</h1>
			<p class="text-gray-600">Convert any website to clean markdown format recursively</p>
		</div>

		<!-- Main Form -->
		<div class="mb-8 rounded-lg bg-white p-6 shadow-lg">
			<form on:submit|preventDefault={handleSubmit} class="space-y-6">
				<!-- URL Input -->
				<div>
					<label for="url" class="mb-2 block text-sm font-medium text-gray-700">
						🌐 Website URL
					</label>
					<input
						id="url"
						type="url"
						bind:value={url}
						placeholder="https://example.com"
						class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-blue-500 focus:ring-2 focus:ring-blue-500"
						disabled={isLoading}
						required
					/>
				</div>

				<!-- Options -->
				<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
					<div>
						<label for="depth" class="mb-2 block text-sm font-medium text-gray-700">
							📊 Max Depth
						</label>
						<select
							id="depth"
							bind:value={maxDepth}
							class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-blue-500 focus:ring-2 focus:ring-blue-500"
							disabled={isLoading}
						>
							<option value={1}>1 level</option>
							<option value={2}>2 levels</option>
							<option value={3}>3 levels</option>
							<option value={4}>4 levels</option>
							<option value={5}>5 levels</option>
						</select>
					</div>

					<div>
						<label for="delay" class="mb-2 block text-sm font-medium text-gray-700">
							⏱️ Delay (ms)
						</label>
						<select
							id="delay"
							bind:value={delay}
							class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-blue-500 focus:ring-2 focus:ring-blue-500"
							disabled={isLoading}
						>
							<option value={500}>500ms</option>
							<option value={1000}>1000ms</option>
							<option value={2000}>2000ms</option>
							<option value={3000}>3000ms</option>
						</select>
					</div>

					<div class="flex items-center">
						<label class="flex items-center space-x-2 pt-6">
							<input
								type="checkbox"
								bind:checked={followExternal}
								class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
								disabled={isLoading}
							/>
							<span class="text-sm font-medium text-gray-700">🌐 Follow external links</span>
						</label>
					</div>
				</div>

				<!-- Submit Button -->
				<div class="flex justify-center">
					<button
						type="submit"
						disabled={isLoading || !url}
						class="rounded-lg bg-blue-600 px-8 py-3 font-medium text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-gray-400"
					>
						{isLoading ? '🔄 Scraping...' : '🚀 Start Scraping'}
					</button>
				</div>
			</form>
		</div>

		<!-- Error Display -->
		{#if error}
			<div class="mb-6 rounded-lg border border-red-200 bg-red-50 p-4">
				<div class="flex items-center">
					<span class="font-medium text-red-600">❌ Error:</span>
					<span class="ml-2 text-red-700">{error}</span>
				</div>
			</div>
		{/if}

		<!-- Results Display -->
		{#if results}
			<div class="rounded-lg bg-white p-6 shadow-lg">
				<!-- Stats -->
				<div class="mb-6">
					<h2 class="mb-4 text-2xl font-bold text-gray-900">✅ Scraping Complete!</h2>
					<div class="mb-4 grid grid-cols-1 gap-4 md:grid-cols-4">
						<div class="rounded-lg bg-blue-50 p-4">
							<div class="text-2xl font-bold text-blue-600">{results.stats?.totalPages || 0}</div>
							<div class="text-sm text-blue-700">Total Pages</div>
						</div>
						<div class="rounded-lg bg-green-50 p-4">
							<div class="text-2xl font-bold text-green-600">
								{results.stats?.successPages || 0}
							</div>
							<div class="text-sm text-green-700">Successful</div>
						</div>
						<div class="rounded-lg bg-red-50 p-4">
							<div class="text-2xl font-bold text-red-600">{results.stats?.errorPages || 0}</div>
							<div class="text-sm text-red-700">Errors</div>
						</div>
						<div class="rounded-lg bg-gray-50 p-4">
							<div class="text-xl font-bold text-gray-600">
								{results.stats?.processingTime || 'N/A'}
							</div>
							<div class="text-sm text-gray-700">Processing Time</div>
						</div>
					</div>

					<!-- Download Buttons -->
					<div class="flex flex-wrap gap-4">
						<button
							on:click={downloadAsMarkdown}
							class="rounded-lg bg-green-600 px-4 py-2 text-white transition-colors hover:bg-green-700"
						>
							📄 Download Markdown
						</button>
						<button
							on:click={downloadAsJson}
							class="rounded-lg bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
						>
							📦 Download JSON
						</button>
					</div>
				</div>

				<!-- Pages Preview -->
				<div class="border-t pt-6">
					<h3 class="mb-4 text-lg font-bold text-gray-900">📄 Pages Preview</h3>
					<div class="max-h-96 space-y-4 overflow-y-auto">
						{#each results.pages as page, index}
							<div class="rounded-lg border border-gray-200 p-4">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<h4 class="font-medium text-gray-900">
											{page.error ? '❌' : '✅'}
											{page.title || page.url}
										</h4>
										<p class="mt-1 text-sm text-gray-600">
											{page.url} (Depth: {page.depth})
										</p>
										{#if page.error}
											<p class="mt-2 text-sm text-red-600">{page.error}</p>
										{:else}
											<p class="mt-2 text-sm text-gray-500">
												{page.markdown?.slice(0, 200)}...
											</p>
										{/if}
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
