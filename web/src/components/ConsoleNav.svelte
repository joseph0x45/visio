<script lang="ts">
	export let user_data: {
		username: string;
		avatar: string;
	};
</script>

<div class=" min-h-full">
	<nav class="bg-white border-b border-gray-200">
		<div class="mx-auto max-w-7xl px-4 lg:px-8">
			<div class="flex h-16 items-center justify-between">
				<div class="flex items-center">
					<a href="/console">
						<h1 class="font-semibold text-lg">VISIO</h1>
					</a>
					<div class="hidden md:block">
						<div class="ml-10 flex items-baseline space-x-4">
							<a
								href="/console/keys"
								class=" rounded-md px-3 py-2 border border-gray-200 text-sm font-medium"
								aria-current="page">Keys</a
							>
							<a href="https://github.com/TheWisePigeon/visio#README" target="_blank" class=" rounded-md px-3 py-2 text-sm font-medium">Docs</a>
							<a
								href="https://github.com/TheWisePigeon/visio"
								target="_blank"
								class=" rounded-md px-3 py-2 text-sm font-medium">GitHub</a
							>
						</div>
					</div>
				</div>
				<div class="hidden md:block">
					<div class="ml-4 flex items-center">
						<span class="">{user_data.username}</span>
						<div class="ml-3">
							<img class="h-8 w-8 rounded-full" src={user_data.avatar} alt="" />
						</div>
					</div>
				</div>
			</div>
		</div>
	</nav>
</div>
