import { redirect } from "@sveltejs/kit"
import type { LayoutServerLoad } from "./$types"

export const load: LayoutServerLoad = async ({ cookies }) => {
  const auth_token = cookies.get("auth_token")
  if (!auth_token) {
    throw redirect(301, "/")
  }
  const response = await fetch(
    "http://localhost:8080/user",
    {
      headers: {
        "Authorization": `Bearer ${auth_token}`
      }
    }
  )
  const data = await response.json() as { avatar: string, email: string, username: string, plan: string }
  return data
}
