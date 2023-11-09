import { redirect, type Actions, fail } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ cookies }) => {
  const auth_token = cookies.get("auth_token")
  if (!auth_token) {
    throw redirect(301, "/")
  }
  const response = await fetch(
    "http://localhost:8080/keys",
    {
      headers: {
        Authorization: `Bearer ${auth_token}`
      }
    }
  )
  const data = await response.json()
  console.log(data)
  return {
    data
  }
}

export const actions : Actions = {
  create: async ({ cookies }) =>{
    try {
      const auth_token = cookies.get("auth_token")
      if (!auth_token){
        throw redirect(301, "/")
      }
      const response = await fetch(
        "http://localhost:8080/keys",
        {
          method:"POST"
        }
      )
      if (response.status==201){
        const { key } = await response.json() as { key: string }
        return {
          key
        }
      }
      return fail(response.status)
    } catch (error) {
      console.log(`Error while creating key: ${error}`)
      return fail(500)
    }
  }
}
