<script lang="ts">
	import Key from '../../../../components/Key.svelte';
	import { type SubmitFunction, enhance } from '$app/forms';
	import toast from 'svelte-french-toast';
	import type { PageData } from './$types';
	export let data: PageData;
	$: ({ keys } = data);
	let show_key = false;
	let created_key = '';
	let loading = false;
	async function copy_key() {
		await navigator.clipboard.writeText(created_key);
		toast.success('Key copied to clipboard');
	}
	function dismiss() {
		show_key = false;
		created_key = '';
	}
	const handle_key_creation: SubmitFunction = () => {
		loading = true;
		toast.loading('Creating your key', { id: 'create_key' });
		return async ({ update, result }) => {
			loading = false;
			toast.dismiss('create_key');
			switch (result.type) {
				case 'success':
					const key = result.data!.key;
					if (!key) {
						toast.error('Something went wrong\nTry again or contact us');
						return;
					}
					created_key = key;
					show_key = true;
					break;
				case 'failure':
					if (result.status == 403) {
						toast.error('You can only have one key currently');
						break;
					}
					toast.error('Something went wrong\nTry again or contact us');
					break;
			}
			await update();
		};
	};
</script>

{#if show_key}
	<div
		class="z-30 fixed inset-0 h-full w-full flex flex-col justify-center items-center bg-black/50"
	>
		<div class="bg-white w-fit rounded-md p-4 flex flex-col gap-5">
			<h1 class="text-center">
				Here is your key! Once you dismiss this modal you won't be able to see it again.
			</h1>
			<div class="p-2 flex justify-center gap-10 border rounded-md">
				<h1>{created_key}</h1>
				<button on:click={copy_key}>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="1.5"
						stroke="currentColor"
						class="w-6 h-6"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184"
						/>
					</svg>
				</button>
			</div>
			<button on:click={dismiss} class="bg-black text-white rounded-md p-2">Dismiss</button>
		</div>
	</div>
{/if}

<div class="min-h-full flex justify-between items-center mb-10">
	<h1 class="text-lg">Your keys ( {keys.length}/1 )</h1>
	<form action="?/create" method="post" use:enhance={handle_key_creation}>
		<button disabled={loading} class="rounded-md p-2 w-fit bg-black text-white">
			<h1>Create a key</h1>
		</button>
	</form>
</div>
{#each keys as key}
	<Key key_prefix={key.prefix} />
{/each}
