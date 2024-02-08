import {QueryClient} from "@tanstack/react-query";
import {Globals} from "./globals";
import {tryGetPrefetch} from "../utils/prefetch";
import {Navigation} from "./navigation";

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
            refetchOnWindowFocus: false,
            staleTime: 5 * 60 * 1000, // 5 minutes
        }
    }
});

function initGlobals() {
    queryClient.setQueryData<Globals>(['globals'], {
        viewerId: tryGetPrefetch('viewerId') || '',
        viewerName: tryGetPrefetch('viewerName') || '',
    }, {})

    queryClient.setQueryData<Navigation>(['navigation'], {
        url: window.location.pathname,
    })
}

initGlobals(); // TODO think about it
