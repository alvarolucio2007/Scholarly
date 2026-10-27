import { PUBLIC_API_URL } from '$env/static/public';

export async function apiFetch(path: string, init: RequestInit = {}): Promise<Response> {
    const token = localStorage.getItem('paseto_token');
    const response = await fetch(`${PUBLIC_API_URL}${path}`, {
        ...init,
        headers: { ...init.headers, Authorization: `Bearer ${token}` }
    });
    if (response.status === 401) {
        localStorage.removeItem('paseto_token');
        window.location.href = '/login';
    }
    return response;
}
