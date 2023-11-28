<script lang="ts">
	import { goto } from '$app/navigation';
	async function github_auth() {
		try {
			let response = await fetch(`http://localhost:8080/auth/request`);
			const { url } = (await response.json()) as { url: string };
			goto(url);
		} catch (err) {
			console.log(err);
		}
	}
</script>

<button on:click={github_auth}>
  <h1>Login with Github</h1>
</button>
