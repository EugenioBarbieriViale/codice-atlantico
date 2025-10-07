<script lang="ts">
	let status: string | null = null;
	let loading = false;

	async function checkHealth() {
		loading = true;
		status = null;
		try {
			const res = await fetch('http://localhost:8080/healthz');
			const data = await res.json();
			status = data.status;
		} catch (err) {
			status = 'unreachable';
		} finally {
			loading = false;
		}
	}
</script>

<main class="h-screen flex flex-col items-center justify-center gap-4">
	<h1 class="text-3xl font-semibold">Codice Atlantico</h1>
	<button
		class="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition disabled:opacity-50"
		on:click={checkHealth}
		disabled={loading}
	>
		{loading ? 'Checking...' : 'Check Backend Health'}
	</button>

	{#if status}
		<p class="text-lg">
			Status:
			<span
				class={
					status === 'ok'
						? 'text-green-600 font-bold'
						: 'text-red-600 font-bold'
				}
			>
				{status}
			</span>
		</p>
	{/if}
</main>

<style>
	main {
		font-family: system-ui, sans-serif;
	}
</style>
