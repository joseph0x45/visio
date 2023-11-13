import { redirect } from "@sveltejs/kit"
import type { LayoutServerLoad } from "./$types"
import { API_URL } from "$lib/config"

export const load: LayoutServerLoad = async ({ cookies }) => {
  console.log("running from server load")
  const auth_token = cookies.get('auth_token')
  console.log(auth_token)
  const response = await fetch(
    `${API_URL}/user`,
    {
      headers: {
        "Authorization": `Bearer ${auth_token}`
      }
    }
  )
  if (response.status!=200){
    console.log(response.status)
    throw redirect(301, "/")
  }
  const data = await response.json() as { avatar: string, email: string, username: string, plan: string }
  return data
}
