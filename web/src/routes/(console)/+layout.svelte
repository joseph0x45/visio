<script lang="ts">
	import type { LayoutData } from './$types';
	import ConsoleNav from '../../components/ConsoleNav.svelte';
	export let data: LayoutData;
	$: ({ username, avatar } = data);
</script>

<div class="">
	<ConsoleNav user_data={{ username, avatar }} />
</div>
