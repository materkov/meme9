export const cookieAuthToken = "authToken";

export function getCookie(name: string): string {
    const value = "; " + document.cookie;
    const parts = value.split("; " + name + "=");
    if (parts.length >= 2) {
        return parts.pop()?.split(";").shift() || "";
    } else {
        return "";
    }
}

export function setCookie(name: string, val: string) {
    document.cookie = name + "=" + val + "; Path=/";
}

export function delCookie(name: string) {
    document.cookie = name + "=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT";
}
