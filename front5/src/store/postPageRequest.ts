import {UseQueryOptions} from "@tanstack/react-query";
import {PostsListRequest, PostsListResponse} from "./types";
import {apiRequest} from "./store";

export const postPageRequest = (id: string): UseQueryOptions<[PostsListResponse]> => ({
    queryKey: ["post", id],
    queryFn: () => apiRequest([
        PostsListRequest({
            fields: "text,user(name)",
            id: [id],
        }),
    ])
})
