import {ResolveRouteRequest, UniversalRenderer} from "./api/renderer";
import {api} from "./Api";

export function resolveRoute(route: string): Promise<UniversalRenderer> {
    return api<ResolveRouteRequest, UniversalRenderer>(
        "meme.Utils.ResolveRoute", {url: route}
    );
}
