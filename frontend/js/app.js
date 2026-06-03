const API_URL = "http://localhost:8080";

async function getMenu() {
    try {
        const response = await fetch(`${API_URL}/menu`)

        if (!response.ok) {
            throw new Error("Failed to fetch menu")
        }

        return await response.json();
    }catch (error) {
        console.error(error);
        return[];
    }
}