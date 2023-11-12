import { redirect, type Actions, fail } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { API_URL } from "$lib/config";

type Key = {
  id: string,
  owner: string,
  prefix: string,
  key_hash: string
}

export const load: PageServerLoad = async ({ cookies }) => {
  const auth_token = cookies.get("auth_token")
  if (!auth_token) {
    throw redirect(301, "/")
  }
  const response = await fetch(
    `${API_URL}/keys`,
    {
      headers: {
        Authorization: `Bearer ${auth_token}`
      }
    }
  )
  const data = await response.json() as Array<Key>
  return {
    keys: data
  }
}

export const actions: Actions = {
  create: async ({ cookies }) => {
    try {
      const auth_token = cookies.get("auth_token")
      if (!auth_token) {
        throw redirect(301, "/")
      }
      const response = await fetch(
        `${API_URL}/keys`,
        {
          method: "POST",
          headers: {
            "Authorization": `Bearer ${auth_token}`
          }
        }
      )
      if (response.status == 201) {
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
  },
  delete: async ({ cookies, request }) => {
    try {
      const auth_token = cookies.get("auth_token")
      if (!auth_token) {
        throw redirect(301, "/")
      }
      const data = await request.formData()
      const key_prefix = data.get('prefix')
      if (!key_prefix) {
        return fail(500)
      }
      const response = await fetch(
        `${API_URL}/keys/${key_prefix}`,
        {
          method: "DELETE",
          headers: {
            "Authorization": `Bearer ${auth_token}`
          }
        }
      )
      if (response.status == 200) {
        return {}
      }
      return fail(response.status)
    } catch (error) {
      console.log(`Error while deleting key ${error}`)
      return fail(500)
    }
  }
}
