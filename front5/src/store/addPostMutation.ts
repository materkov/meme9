import {FeedListResponse, Post} from "./types";
import {feedRequest} from "./feedRequest";
import {queryClient, updateQuery} from "./store";
import {UseMutationOptions} from "@tanstack/react-query";

export const addPostMutation = (text: string): UseMutationOptions => ({
    mutationFn: () => {
        return new Promise<void>((resolve, reject) => {
            setTimeout(() => {
                resolve();

                const post: Post = {
                    id: Math.random().toString(),
                    text: text,
                    date: Math.floor((new Date()).getTime()),
                    user: {
                        id: "user10",
                        name: "User 10 name",
                    },
                }

                updateQuery(feedRequest().queryKey, (data: [FeedListResponse]) => {
                    data[0].items = [post, ...data[0].items];
                });
            }, 1000);
        })
    }
})
