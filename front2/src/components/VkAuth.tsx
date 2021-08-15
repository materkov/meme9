import React from "react";
import {VkAuthRenderer} from "../types";

export const VkAuth = (props: { data: VkAuthRenderer }) => {
    return <a href={props.data.url}>Войти через VK</a>;
}
