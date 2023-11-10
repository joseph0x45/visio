<script lang="ts">
	import { goto } from '$app/navigation';
	import toast from 'svelte-french-toast';
	let loading = false;
	async function github_auth() {
		try {
			loading = true;
			toast.loading('Signing you in', {
				id: 'loading'
			});
			let response = await fetch('http://localhost:8080/auth/request');
			const { url } = (await response.json()) as { url: string };
			goto(url);
		} catch (err) {
			loading = false;
			toast.dismiss('loading');
			toast.error('Something went wrong. Please retry or contact us');
			console.log(err);
		}
	}
</script>

<title>Visio | Landing</title>

<div class="bg-white">
	<header class="absolute inset-x-0 top-0 z-50">
		<nav class="flex items-center justify-between p-6 lg:px-8" aria-label="Global">
			<div class="flex lg:flex-1">
				<a href="/" class="-m-1.5 p-1.5">
					<span class="sr-only">Visio</span>
					<img
						class="h-12 w-auto"
						src="/logo.png"
						alt="logo"
					/>
				</a>
			</div>
			<div class="hidden lg:flex lg:gap-x-12">
				<a href="https://github.com/thewisepigeon/visio" target="_blank" class="text-sm font-semibold leading-6 text-gray-900">GitHub</a>
				<a href="/" class="text-sm font-semibold leading-6 text-gray-900">Documentation</a>
				<a href="https://discord.gg/rDNvUpct" target="_blank" class="text-sm font-semibold leading-6 text-gray-900">Discord</a>
			</div>
			<button on:click={github_auth} class="hidden lg:flex lg:flex-1 lg:justify-end">
				<h1 class="text-sm font-semibold leading-6 text-gray-900">
          Log in with GitHub <span aria-hidden="true">&rarr;</span>
        </h1>
			</button>
		</nav>
	</header>

	<div class="relative isolate px-6 pt-14 lg:px-8">
		<div
			class="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80"
			aria-hidden="true"
		>
			<div
				class="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]"
				style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
			/>
		</div>
		<div class="mx-auto max-w-2xl py-32 sm:py-48 lg:py-56">
			<div class="hidden sm:mb-8 sm:flex sm:justify-center">
				<div
					class="relative rounded-full px-3 py-1 text-sm leading-6 text-gray-600 ring-1 ring-gray-900/10 hover:ring-gray-900/20"
				>
					<h1>Visio is still in active development. Join our <a  class="text-sm font-semibold leading-6 text-gray-900" href="https://discord.gg/rDNvUpct" target="_blank">Discord</a> to stay up to date and contribute</h1>
				</div>
			</div>
			<div class="text-center">
				<h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
					Supercharge your apps with face detection and recognition
				</h1>
				<p class="mt-6 text-lg leading-8 text-gray-600">
					Visio is an open source cloud based service that provides algorithms for face detection
					and recognition.
				</p>
				<div class="mt-10 flex items-center justify-center gap-x-6">
					<a
						href="/#"
						class="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
						>Get started</a
					>
					<a href="/#" class="text-sm font-semibold leading-6 text-gray-900"
						>Learn more <span aria-hidden="true">→</span></a
					>
				</div>
			</div>
		</div>
		<div
			class="absolute inset-x-0 top-[calc(100%-13rem)] -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[calc(100%-30rem)]"
			aria-hidden="true"
		>
			<div
				class="relative left-[calc(50%+3rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%+36rem)] sm:w-[72.1875rem]"
				style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
			/>
		</div>
	</div>
</div>
