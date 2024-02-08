import {useQuery} from "@tanstack/react-query";
import {queryClient} from "./queryClient";

export interface Navigation {
    url: string;
}

export function useNavigation(): Navigation {
    const data = useQuery<Navigation>({queryKey: ['navigation']});
    if (!data.data) {
        throw 'No global navigation';
    }

    return data.data;
}

export function navigationGo(url: string) {
    window.history.pushState(null, '', url);

    queryClient.setQueryData<Navigation>(['navigation'], {
        url: url,
    })
}
