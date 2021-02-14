<template>
  <div id="app">
    <div class="page-container">
      <md-app md-mode="reveal">
        <md-app-toolbar class="md-primary">
          <md-button class="md-icon-button" @click="menuVisible = !menuVisible">
            <md-icon>menu</md-icon>
          </md-button>
          <span class="md-title" style="flex: 1">cfd</span>
          <!-- <md-button>Refresh</md-button>
          <md-button>Create</md-button> -->
        </md-app-toolbar>
        <md-app-drawer :md-active.sync="menuVisible">
          <md-toolbar class="md-transparent" md-elevation="0">
            <span class="md-title">Certificate Authorities</span>
          </md-toolbar>

          <md-list>
            <md-list-item
              v-for="(id, index) in caIds"
              :key="index"
              @click="
                mutate({ property: 'caId', with: id });
                menuVisible = false;
              "
            >
              <md-icon>segment</md-icon>
              <span class="md-list-item-text">{{ id }}</span>
            </md-list-item>

            <md-divider></md-divider>

            <md-list-item>
              <md-icon>library_add_check</md-icon>
              <span class="md-list-item-text"
                >Add an existing CA to the list</span
              >
            </md-list-item>

            <md-list-item>
              <md-icon>library_add</md-icon>
              <span class="md-list-item-text">Create new CA</span>
            </md-list-item>

            <md-divider></md-divider>

            <md-list-item>
              <span class="md-list-item-text" style="align-items: center"
                >version: {{ version }}</span
              >
            </md-list-item>
          </md-list>
        </md-app-drawer>

        <md-app-content>
          <certificates></certificates>
        </md-app-content>
      </md-app>
    </div>
  </div>
</template>

<script>
import { mapState, mapMutations } from "vuex";
import Certificates from "./components/Certificates.vue";

export default {
  name: "App",
  components: {
    Certificates,
  },
  data: () => ({
    menuVisible: false,
  }),
  computed: {
    ...mapState(["version", "caIds", "caId"]),
  },
  methods: {
    fetchStatus() {
      this.$store.dispatch("fetchData", {
        method: "get",
        url: "status",
        key: "version",
        subkey: "version",
      });
    },
    ...mapMutations(["mutate"]),
  },
  watch: {
    caId(newValue, oldValue) {
      if (newValue == oldValue) {
        return;
      }
      if (newValue != undefined || newValue != "") {
        this.$store.dispatch("fetchData", {
          method: "get",
          url: `v1/ca/${newValue}/certificates?parse=true`,
          key: "certs",
          subkey: "",
        });
      }
    },
  },
  mounted() {
    this.fetchStatus();
  },
};
</script>

<style lang="scss" scoped>
// Demo purposes only
.md-app {
  min-height: 100vh;
}
.md-drawer {
  width: 350px;
  max-width: calc(100vw - 125px);
  min-height: 100%;
}
</style>