<script lang="ts">
	import { type SubmitFunction, enhance } from '$app/forms';
	import toast from 'svelte-french-toast';

	export let key_prefix: string;
	const handle_key_deletion: SubmitFunction = async ({ data }) => {
		toast.loading('Revoking key', { id: 'revoke_key' });
		console.log(data.get('prefix'));
		toast.dismiss('revoke_key');
		return async ({ update, result }) => {
			switch (result.type) {
				case 'failure':
					toast.error('Something went wrong\nTry again or contact us');
					break;
				case 'success':
					toast.success('Key deleted');
					break;
				default:
					break;
			}
      await update()
		};
	};
</script>

<div class="rounded-md w-full p-2 h-16 border border-gray-400 flex items-center justify-between">
	<h1>{key_prefix}**************************</h1>
	<form action="?/delete" method="post" use:enhance={handle_key_deletion}>
		<input name="prefix" type="text" hidden value={key_prefix} />
		<button class=" ">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="w-6 h-6 hover:text-red-500"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"
				/>
			</svg>
		</button>
	</form>
</div>
