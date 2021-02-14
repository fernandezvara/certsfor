<template>
  <div v-if="certs.length == 0">
    <md-empty-state
      md-icon="security"
      md-label="Create a Certificate Authority"
      md-description="With a CA you'll be able to create the certificates to empower your systems and applications."
    >
      <md-button class="md-primary md-raised">Create your first CA</md-button>
    </md-empty-state>
  </div>
  <div v-else>
    <certificate v-for="(cert, index) in certs" :cert="cert," :key="index" />
  </div>
</template>

<script>
import { mapState, mapActions } from "vuex";
import Certificate from "./Certificate";

export default {
  components: {
    Certificate,
  },
  computed: {
    ...mapState(["certs"]),
  },
  methods: {
    ...mapActions(["fetchData"]),
  },
  mounted() {
    if (this.$store.ca_id == undefined || this.$store.ca_id == "") {
      return;
    }
    this.fetchData({
      method: "get",
      url: `v1/ca/${this.$store.ca_id}/certificates?parse=true`,
      key: "certs",
    });
  },
};
</script>

<style lang="scss" scoped>
</style>
