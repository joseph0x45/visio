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
  try {
    const response = await fetch(`${API_URL}/auth/user`, {
      credentials: "include",
      headers: {
        "X-VISIO-APP-IDENTIFIER": "app_identifier"
      }
    })
    if (response.status == 200) {
      const userData = await response.json() as UserData
      console.log(userData)
      return {
        data: userData
      }
    }
    console.log(response.status)
    return redirect(302, "/?error=goesbrrrrr")
  } catch (error) {
    console.log("Error while fetching user data: ")
    console.log(error)
    throw redirect(302, "/?context=internal")
  }
}
