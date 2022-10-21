import {QueryClient, useQuery as useQueryReact, UseQueryOptions} from "@tanstack/react-query";

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: 15 * 60 * 1000,
        }
    }
});

export function fetcher({queryKey}: any): Promise<any> {
    return new Promise((resolve, reject) => {
        const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me/api" : "http://localhost:8000/api";

        fetch(apiHost + "?urls=" + queryKey[0], {
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken')
            }
        })
            .then(r => r.json())
            .then(r => {
                let result = null;
                for (let item of r) {
                    queryClient.setQueryData([item.url], item, {});

                    if (item.url === queryKey[0]) {
                        result = item;
                    }
                }

                resolve(result);
            })
    })
}

export function useQuery<T>(url: string) {
    return useQueryReact<T>([url], fetcher);
}
