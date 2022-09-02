import {UseQueryOptions} from "@tanstack/react-query";
import {FeedList, FeedListResponse, PostsListResponse, UsersListResponse} from "./types";
import {apiRequest, queryClient} from "./store";
import {profileRequest} from "./profileRequest";
import {postPageRequest} from "./postPageRequest";

export const feedRequest = (): UseQueryOptions<[FeedListResponse]> => ({
    queryKey: ["feed"],
    onSuccess: (data) => {
        for (let post of data[0].items) {
            if (post.user?.id) {
                const userQueryKey = profileRequest(post.user?.id).queryKey;
                if (userQueryKey) {
                    const data: UsersListResponse[] = [{
                        items: [post.user],
                    }];
                    queryClient.setQueryData(userQueryKey, data);
                }
            }

            const postQueryKey = postPageRequest(post.id).queryKey;
            if (postQueryKey) {
                const data: PostsListResponse[] = [{
                    items: [post],
                }];

                queryClient.setQueryData(postQueryKey, data);
            }
        }
    },
    queryFn: () => apiRequest([
        FeedList({
            fields: "foo,items(text,user(id,name,posts(id,text,user(id,name))))",
        }),
    ]),
})
