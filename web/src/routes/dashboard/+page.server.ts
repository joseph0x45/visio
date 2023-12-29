import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { API_URL } from "$lib/config";

type UserData = {
  id: string,
  email: string,
  avatar: string,
  username: string
}

export const load: PageServerLoad = async ({ cookies }) => {
  const token = cookies.get("token") ?? ""
  if (token == "") {
    throw redirect(302, "/?error=no_auth_cookie")
  }
  try {
    const response = await fetch(`${API_URL}/auth/user`, {
      headers: {
        "X-VISIO-APP-IDENTIFIER": "app_identifier",
        "Authorization": `Bearer ${token}`
      }
    })
    if (response.status == 200) {
      const userData = await response.json() as UserData
      console.log(userData)
      return {
        data: userData
      }
    }
    console.log(`Got HTTP ${response.status} while requesting for user data`)
    throw redirect(302, `/?error=http_${response.status}`)
  } catch (error) {
    console.log("Error while making HTTP request ", error)
    throw redirect(302, "/?error=internal")
  }
}
