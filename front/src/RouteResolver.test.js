import {test} from "@jest/globals";
import {resolveRoute} from "./RouteResolver";

test('resolve route', () => {
    resolveRoute('/test').then(r => {
        expect(r).toBeDefined();
    })
});
