<script lang="ts">
	import '../app.css';
	import { Toaster } from 'svelte-french-toast';
	import { navigating } from '$app/stores';
	import { browser } from '$app/environment';
	import NProgress from 'nprogress';
	import 'nprogress/nprogress.css';
	NProgress.configure({
		minimum: 0.16
	});
	$: {
		if (browser) {
			if ($navigating) {
				NProgress.start();
			}
			if (!$navigating) {
				NProgress.done();
			}
		}
	}
</script>

<Toaster />
<slot />
