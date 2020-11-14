import React from "react";
import * as schema from "../../schema/login";
import {Header} from "../Header/Header";
import {PostComposer} from "./PostComposer";
import {UnathorizedPlaceholder} from "./UnathorizedPlaceholder";

export interface ComposerProps {
    data: schema.ComposerRenderer;
}

export function Composer(props: ComposerProps) {
    return (
        <div>
            <Header data={props.data.headerRenderer}/>

            {props.data.unathorizedText ?
                <UnathorizedPlaceholder data={props.data}/> :
                <PostComposer data={props.data}/>
            }
        </div>
    );
}
