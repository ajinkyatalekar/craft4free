const API_URL = import.meta.env.PROD
  ? "http://api.craft4free.online/"
  : "http://api.craft4free.online/";

export async function createServer(server: string, access_token: string) {
  try {
    const response = await fetch(`${API_URL}servers?server=${server}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${access_token}`,
      },
    });

    return response.json();
  } catch (error) {
    console.error("API request failed:", error);
    throw error;
  }
}
