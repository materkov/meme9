import {QueryClient, QueryKey} from "@tanstack/react-query";

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            refetchOnWindowFocus: false,
            staleTime: 60 * 1000,
        }
    }
});

// @ts-ignore
window.q = queryClient;

type Operation = {
    method: string;
    params: any;
}

export function apiRequest(operations: Operation[]): Promise<any> {
    return fetch("http://127.0.0.1:8000/api", {
        method: 'POST',
        body: JSON.stringify(operations),
    })
        .then(r => r.json());
}

export function updateQuery(key?: QueryKey, updater?: (data: any) => any) {
    if (!key || !updater) {
        return;
    }

    let data = queryClient.getQueryData(key);

    // TODO ImmerJs?
    if (data) {
        data = JSON.parse(JSON.stringify(data));
    }

    updater(data);

    queryClient.setQueryData(key, data);
}
