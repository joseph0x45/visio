import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ cookies, url, request }) => {
  const session = url.searchParams.get("session") ?? ""
  const appIdentifier = request.headers.get("X-APP-IDENTIFIER") ?? ""
  if (appIdentifier != "app_identifier") {
    throw redirect(302, "/")
  }
  if (session == "") {
    throw redirect(302, "/")
  }
  cookies.set('session', session, {
    path: "/",
    sameSite: "lax",
    secure: true
  })
  throw redirect(302, "/?hello=bozo")
}
