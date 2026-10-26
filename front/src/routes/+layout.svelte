<script lang="ts">
	import './layout.css';
	import favicon from '#lib/assets/favicon.svg';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<svelte:head>
    <title>Scholarly</title>
    <meta name="description" content="Sistema de Gestão Acadêmica" />
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
