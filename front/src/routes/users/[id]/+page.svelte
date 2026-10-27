<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/state'; // Atualizado para Svelte 5
    import { resolve } from '$app/paths';

    // Pega o ID da URL de forma reativa
    let userId = $derived(page.params.id);

    let name = $state('');
    let email = $state('');
    let cpf = $state('');
    let password = $state('');

    let isLoading = $state(true);
    let isSaving = $state(false);
    let errorMessage = $state('');
    let successMessage = $state('');

    function maskCPF(value: string) {
        const digits = value.replace(/\D/g, '');
        const limited = digits.slice(0, 11);
        if (limited.length <= 3) return limited;
        if (limited.length <= 6) return `${limited.slice(0, 3)}.${limited.slice(3)}`;
        if (limited.length <= 9) return `${limited.slice(0, 3)}.${limited.slice(3, 6)}.${limited.slice(6)}`;
        return `${limited.slice(0, 3)}.${limited.slice(3, 6)}.${limited.slice(6, 9)}-${limited.slice(9)}`;
    }

    function handleCPFInput(event: Event) {
        const target = event.target as HTMLInputElement;
        cpf = maskCPF(target.value);
    }

    onMount(async () => {
        const token = localStorage.getItem('paseto_token');
        if (!token) {
            window.location.href = '/login';
            return;
        }

        try {
            const response = await fetch(`http://localhost:8080/users/${userId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (!response.ok) throw new Error('User not found');

            const data = await response.json();
            const user = data.data ? data.data : data;

            name = user.name;
            email = user.email;
            cpf = maskCPF(user.cpf || '');
            password = user.password;

        } catch (err) {
            errorMessage = 'Failed to load user data.';
        } finally {
            isLoading = false;
        }
    });

    async function handleUpdate(event: SubmitEvent) {
        event.preventDefault();
        isSaving = true;
        errorMessage = '';
        successMessage = '';

        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch('http://localhost:8080/users', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    id: Number(userId),
                    name: name.trim(),
                    email: email.trim(), // <-- Remove espaços extras
                    cpf: cpf.replace(/\D/g, ''), // <-- Envia o CPF SÓ COM NÚMEROS
                    // Se o usuário não digitou senha nova, não mande.
                    // (O Go com nilIfEmpty vai ignorar string vazia)
                    password: password.trim() || undefined
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                // Adicione um console.log para ver o erro exato do Go no navegador
                console.log("Erro do Go:", errorData.error);
                throw new Error(errorData.error || 'Failed to update user');
            }

            successMessage = 'User updated successfully!';
            setTimeout(() => { window.location.href = '/users'; }, 1500);

        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:head>
    <title>Edit User | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">

        <a href={resolve('/users')} class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Users
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Edit User #{userId}</h1>

            {#if isLoading}
                <p class="text-gray-500">Loading user data...</p>
            {:else}
                <form onsubmit={handleUpdate} class="flex flex-col gap-5">
                    <div>
                        <label for="edit-name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
                        <input
                            id="edit-name"
                            type="text"
                            bind:value={name}
                            required
                            class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                        />
                    </div>

                    <div>
                        <label for="edit-email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                        <input
                            id="edit-email"
                            type="email"
                            bind:value={email}
                            required
                            class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                        />
                    </div>

                    <div>
                        <label for="edit-cpf" class="block text-sm font-medium text-gray-700 mb-1">CPF</label>
                        <input
                            id="edit-cpf"
                            type="text"
                            value={cpf}
                            oninput={handleCPFInput}
                            required
                            maxlength="14"
                            class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                        />
                    </div>
                    <div>
                        <label for="edit-password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
                        <input
                            id="edit-password"
                            type="text"
                            bind:value={password}
                            required
                            class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                        />
                    </div>

                    {#if errorMessage}
                        <div class="bg-red-50 text-red-500 text-sm p-3 rounded-lg text-center">{errorMessage}</div>
                    {/if}

                    {#if successMessage}
                        <div class="bg-green-50 text-green-600 text-sm p-3 rounded-lg text-center">{successMessage}</div>
                    {/if}

                    <button
                        type="submit"
                        disabled={isSaving}
                        class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg mt-2 transition-colors"
                    >
                        {isSaving ? 'Saving...' : 'Save Changes'}
                    </button>
                </form>
            {/if}
        </div>
    </div>
</div>
