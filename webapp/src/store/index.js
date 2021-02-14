import axios from 'axios'
import Vue from 'vue'
import Vuex from 'vuex'
Vue.use(Vuex)

// import state from './state'
// import mutations from './mutations'
// import actions from './actions'
// import modules from './modules'

// export default new Vuex.Store({
//   state,
//   mutations,
//   actions,
//   modules
// })

export default new Vuex.Store({
    state: {
        caId: '',
        caIds: [
            '83c99c5d-16ea-44f0-98c8-4ed4d1d1c177'
        ],
        certs: [],
        status: {},
        version: ''
    },
    mutations: {
        mutate(state, payload) {
            state[payload.property] = payload.with;
        }
    },
    actions: {
        // payload
        // {
        //    method: 'get/post'
        //    url: ''
        //    body: {} -> json with the data to send to the API
        //    key: is the property (key) to fill with data
        //    subkey: key to get information (response.data[key]) or empty if you want all the response.data
        // }
        async fetchData({ commit }, payload) {

            let query = {
                method: payload.method,
                url: `http://192.168.1.159:8080/${payload.url}`,
            }

            if (payload.body) {
                query.data = payload.body
            }

            let response = await axios(query)
            if (payload.key) {
                commit('mutate', {
                    property: payload.key,
                    with: payload.subkey == '' ? response.data : response.data[payload.subkey]
                })
            }

            return response.data

        },
    },
    getters: {
        // version(state) {
        //     return state.version
        // },
        // caId(state) {
        //     return state.caId
        // },
        // caIds(state) {
        //     return state.caIds
        // },
        // certs(state) {
        //     return state.certs
        // }
    },
    modules: {

    }
})