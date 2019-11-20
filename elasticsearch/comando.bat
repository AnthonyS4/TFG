curl -X PUT "localhost:9200/_template/template_1?pretty" -H 'Content-Type: application/json' -d'
{
 "template": "packets-*",
  "mappings": {
    "pcap_file": {
      "dynamic": "false",
      "properties": {
        "timestamp": {
          "type": "date"
        },
        "layers": {
          "properties": {
            "frame": {
              "properties": {
                "frame_frame_len": {
                  "type": "long"
                },
                "frame_frame_protocols": {
                  "type": "keyword"
                }
              }
            },
            "ip": {
              "properties": {
                "ip_ip_src": {
                  "type": "ip"
                },
                "ip_ip_dst": {
                  "type": "ip"
                }
              }
            },
            "udp": {
              "properties": {
                "udp_udp_srcport": {
                  "type": "integer"
                },
                "udp_udp_dstport": {
                  "type": "integer"
                }
              }
            }
          }
        }
      }
    }
  }
}
'
