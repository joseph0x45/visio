<script lang="ts">
	import { goto } from '$app/navigation';
	import toast from 'svelte-french-toast';
	import { API_URL } from '$lib/config';
	async function githubAuth() {
		try {
			const response = await fetch(`${API_URL}/auth/url`);
			if (response.status == 200) {
				const { url } = (await response.json()) as { url: string };
				goto(url);
				return;
			}
			console.log(`Error while fetching auth link: Expected HTTP 200 got ${response.status}`);
			toast.error('Something went wrong! Please try again');
		} catch (err) {
			console.log(`Error while fetching auth link: ${err}`);
			toast.error('Something went wrong! Please try again');
		}
	}
</script>

<h1>Welcome to visio</h1>

<button class="p-2 bg-gray-500 text-white rounded-md m-2" on:click={githubAuth}>Login with Github</button>
