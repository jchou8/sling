export default {
    "roots": [
        "./src"
    ],
    "moduleFileExtensions": [
        "ts",
        "tsx",
        "js",
        "jsx"
    ],
    "transform": {
        "\\.(ts|tsx)$": "<rootDir>/node_modules/ts-jest/preprocessor.js"
    },
    "setupFiles": [
        "raf/polyfill"
    ],
    "testRegex": "/__tests__/.*\\.(ts|tsx|js)$",
    "snapshotSerializers": [
        "enzyme-to-json"
    ]
}