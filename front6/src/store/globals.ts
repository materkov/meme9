import {AuthResp} from "../api/api";
import {useQuery} from "@tanstack/react-query";
import {queryClient} from "./queryClient";

export interface Globals {
    viewerId: string;
    viewerName: string;
}

export function useGlobals(): Globals {
    const data = useQuery<Globals>({
        queryKey: ['globals'],
        staleTime: Infinity,
    });
    if (!data.data) {
        throw 'No globals';
    }
    return data.data;
}

export function setAuth(authResp: AuthResp) {
    queryClient.setQueryData<Globals>(['globals'], {
        viewerId: authResp.userId,
        viewerName: authResp.userName,
    })
}