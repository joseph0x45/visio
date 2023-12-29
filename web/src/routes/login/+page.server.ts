import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ cookies, url }) => {
  const token = url.searchParams.get("token") ?? ""
  if (token == "") {
    throw redirect(302, "/")
  }
  cookies.set('token', token, {
    path: "/",
    sameSite: "lax",
    secure: true
  })
  throw redirect(302, "/dashboard")
}
