import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ cookies, url }) => {
  const session = url.searchParams.get("session") ?? ""
  if (session == "") {
    throw redirect(302, "/")
  }
  // cookies.set('auth_token', token, {
  //   path: "/",
  //   sameSite: "lax",
  //   secure: true
  // })
  throw redirect(302, "/?hello=bozo")
}
