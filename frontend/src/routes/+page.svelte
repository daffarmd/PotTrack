<script>
	import { onMount } from 'svelte';
	let message = 'Click the button to ping backend';
	async function ping() {
		try {
			const res = await fetch('/api/health');
			if (!res.ok) throw new Error(res.statusText);
			const data = await res.json();
			message = data.message ?? JSON.stringify(data);
		} catch (err) {
			message = `Error: ${err.message}`;
		}
	}

	onMount(() => {
		// optional auto-ping
	});
</script>

<main class="min-h-screen flex items-center justify-center p-6">
	<div class="max-w-xl w-full bg-white rounded-lg shadow-md p-8">
		<h1 class="text-2xl font-semibold mb-4">PotTrack</h1>
		<p class="mb-4 text-slate-600">A small Svelte + Tailwind frontend integrated with your backend.</p>
		<div class="flex gap-2">
			<button class="px-4 py-2 bg-sky-600 text-white rounded hover:bg-sky-700" on:click={ping}>Ping backend</button>
			<button class="px-4 py-2 border rounded text-slate-700" on:click={() => message = 'Click the button to ping backend'}>Reset</button>
		</div>
		<div class="mt-4 p-4 bg-slate-50 rounded border text-slate-700">{message}</div>
	</div>
</main>
