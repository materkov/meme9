import {QueryClient, useQuery as useQueryReact} from "@tanstack/react-query";
import DataLoader from "dataloader";

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: 15 * 60 * 1000,
        }
    }
});

function loadResources(resources: any): Promise<any> {
    return new Promise((resolve, reject) => {
        const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me/api" : "http://localhost:8000/api";

        fetch(apiHost, {
            method: 'POST',
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken')
            },
            body: JSON.stringify({
                resources: resources,
            })
        })
            .then(r => r.json())
            .then(r => {
                let result: { [key: string]: any } = {};
                for (let item of r) {
                    queryClient.setQueryData([item.url], item, {});
                    result[item.url] = item;
                }

                const t = resources.map((resource: any) => result[resource] ?? null);
                resolve(t);
            })
    })
}

const userLoader = new DataLoader(resources => loadResources(resources), {
    cache: false,
})

export function fetcher({queryKey}: any): Promise<any> {
    return userLoader.load(queryKey[0]);
}

export function useQuery<T>(url: string) {
    return useQueryReact<T>([url], fetcher);
}
