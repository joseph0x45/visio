<script lang="ts">
	import type { LayoutData } from './$types';
	import ConsoleNav from '../../components/ConsoleNav.svelte';
	export let data: LayoutData;
	$: ({ username, avatar } = data);
</script>

<div class="min-h-full bg-gray-100">
	<ConsoleNav user_data={{ username, avatar }} />
	<main class="h-full">
		<div class="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
			<slot />
		</div>
	</main>
</div>
