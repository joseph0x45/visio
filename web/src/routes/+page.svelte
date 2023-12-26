<script lang="ts">
	import { goto } from '$app/navigation';
	function githubAuth() {
		let client_id = 'dc0bef6ce892da0d4d39';
		let redirect_uri = 'http://localhost:8080/auth/callback';
		let githubAuthUrl = `https://github.com/login/oauth/authorize?client_id=${client_id}&redirect_uri=${redirect_uri}&scope=user:email`;
		goto(githubAuthUrl);
	}
</script>

<h1>Welcome to visio</h1>

<button on:click={githubAuth}>Login with Github</button>
