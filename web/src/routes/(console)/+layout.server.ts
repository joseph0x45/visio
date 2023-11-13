import { redirect } from "@sveltejs/kit"
import type { LayoutServerLoad } from "./$types"
import { API_URL } from "$lib/config"

export const load: LayoutServerLoad = async ({ cookies, request }) => {
  const auth_token = cookies.get("auth_token")
  if (!auth_token) {
    throw redirect(301, "/")
  }
  const response = await fetch(
    `${API_URL}/user`,
    {
      headers: {
        "Authorization": `Bearer ${auth_token}`
      }
    }
  )
  const data = await response.json() as { avatar: string, email: string, username: string, plan: string }
  return data
}
