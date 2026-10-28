<script lang="ts">
    import { resolve } from '$app/paths';
    import { PUBLIC_API_URL } from '$env/static/public';

    let name = $state('');
    let cpf = $state('');
    let email = $state('');
    let password = $state('');
    let errorMessage = $state('');
    let successMessage = $state('');

    async function handleSignup(event: SubmitEvent) {
        event.preventDefault();
        try {
            const response = await fetch(`${PUBLIC_API_URL}/authentication/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name, cpf, email, password })
            });

            if (!response.ok) {
                const errorData = await response.json();
                errorMessage = errorData.error || 'An error has occurred.';
                successMessage = '';
            } else {
                successMessage = 'User has been created successfully! Please return to the login screen below.';
                errorMessage = '';
            }
        } catch (error) {
            errorMessage = 'Connection error. Please try again.';
            console.log(error);
        }
    }

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
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100 p-4">
    <div class="w-full max-w-md bg-white p-8 rounded-xl shadow-lg">
        <h1 class="text-2xl font-bold text-center text-gray-800 mb-6">Registration</h1>

        <form onsubmit={handleSignup} class="flex flex-col gap-5">
            <div>
                <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
                <input
                    id="name"
                    type="text"
                    bind:value={name}
                    required
                    placeholder="John Doe"
                    class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                />
            </div>
            <div>
                <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                <input
                    id="email"
                    type="email"
                    bind:value={email}
                    required
                    placeholder="johndoe@example.com"
                    class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                />
            </div>
            <div>
                <label for="cpf" class="block text-sm font-medium text-gray-700 mb-1">CPF</label>
                <input
                    id="cpf"
                    type="text"
                    value={cpf}
                    oninput={handleCPFInput}
                    required
                    placeholder="000.000.000-00"
                    maxlength="14"
                    class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                />
            </div>
            <div>
                <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
                <input
                    id="password"
                    type="password"
                    bind:value={password}
                    required
                    placeholder="••••••••"
                    minlength="8"
                    class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                />
            </div>

            {#if errorMessage}
                <div class="bg-red-50 text-red-500 text-sm p-3 rounded-lg text-center">
                    {errorMessage}
                </div>
            {/if}
            {#if successMessage}
                <div class="bg-green-50 text-green-500 text-sm p-3 rounded-lg text-center">
                    {successMessage}
                </div>
            {/if}

            <button
                type="submit"
                class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2.5 px-4 rounded-lg mt-2 transition-colors duration-200"
            >
                Create Account
            </button>
        </form>

        <div class="mt-6 text-center text-sm text-gray-600">
            Already have an account?
            <a href={resolve('/login')} class="text-blue-600 hover:text-blue-800 font-semibold hover:underline transition-all">
                Sign in
            </a>
        </div>
    </div>
</div>
