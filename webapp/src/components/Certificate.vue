<template>
  <md-card md-elevation-1>
    <md-card-header>
      <div class="md-title" :class="{ danger: expired(cert.parsed.not_after) }">
        {{ cert.parsed.dn }}
        <md-badge
          class="md-square"
          md-content="CA"
          style="margin-right: 20px"
          v-if="cert.parsed.is_ca"
        />
      </div>
      <div class="md-subhead">
        Key Usages: <b>{{ cert.parsed.key_usage.join(", ") }}</b
        ><template v-if="cert.parsed.ext_key_usage.length > 0">
          | Extended:</template
        >
        <b>{{ cert.parsed.ext_key_usage.join(", ") }}</b>
      </div>
    </md-card-header>

    <md-card-content>
      <span>
        <b>
          <template
            v-if="expired(cert.parsed.not_after)"
            :class="{ danger: expired(cert.parsed.not_after) }"
            >Expired!!</template
          >
          <template v-else
            >Expires
            {{ formatDistanceUnix(cert.parsed.not_after, now) }}</template
          >
        </b>
        ({{ fromUnixTime(cert.parsed.not_after) }}) </span
      ><br />
      <span
        v-if="
          cert.parsed.dns_names.length +
            cert.parsed.emails.length +
            cert.parsed.ips.length +
            cert.parsed.uris.length >
          0
        "
      >
        SANS:
        <span
          class="label green"
          v-for="(n, i) in cert.parsed.dns_names"
          :key="i"
          >dns <b>{{ n }}</b></span
        >
        <span class="label blue" v-for="(n, i) in cert.parsed.emails" :key="i"
          >email <b>{{ n }}</b></span
        >
        <span class="label orange" v-for="(n, i) in cert.parsed.ips" :key="i"
          >ip <b>{{ n }}</b></span
        >
        <span class="label purple" v-for="(n, i) in cert.parsed.uris" :key="i"
          >uri <b>{{ n }}</b></span
        >
      </span>
    </md-card-content>

    <md-card-actions>
      <md-button v-clipboard="decode64(cert.certificate)"
        >COPY CERTIFICATE</md-button
      >
      <md-button v-clipboard="decode64(cert.key)">COPY KEY</md-button>
    </md-card-actions>
  </md-card>
</template>

<script>
import base64 from "base-64";
import fromUnixTime from "date-fns/fromUnixTime";
import getUnixTime from "date-fns/getUnixTime";
import formatDistance from "date-fns/formatDistance";
import format from "date-fns/format";

export default {
  data() {
    return {
      now: 0,
    };
  },
  props: {
    cert: Object,
  },
  methods: {
    decode64(input) {
      return base64.decode(input);
    },
    expired(date) {
      return this.now > fromUnixTime(date);
    },
    fromUnixTime(date) {
      return format(fromUnixTime(date), "dd-MMM-yyyy HH:mm");
    },
    formatDistanceUnix(date1, date2) {
      return formatDistance(fromUnixTime(date1), fromUnixTime(date2), {
        addSuffix: true,
      });
    },
  },
  created() {
    this.now = getUnixTime(Date.now());
  },
};
</script>

<style lang="scss" scoped>
.md-card {
  margin-top: 10px;
  vertical-align: top;
}
.danger {
  color: #f00;
}
.label {
  padding: 5px;
  margin-right: 5px;
}
.green {
  background-color: rgb(193, 236, 136);
}
.blue {
  background-color: rgb(177, 208, 245);
}
.orange {
  background-color: rgb(250, 229, 160);
}
.purple {
  background-color: rgb(225, 184, 241);
}
</style>