import {apiRequest} from "./store";
import {UsersList, UsersListResponse} from "./types";
import {UseQueryOptions} from "@tanstack/react-query";

export const profileRequest = (id: string): UseQueryOptions<[UsersListResponse]> => ({
    queryKey: ["user", id],
    queryFn: () => apiRequest([
        UsersList({
            id: [id],
            fields: "name"
        })
    ]),
})