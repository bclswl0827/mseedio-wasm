## Exported API

The allowlist contains 88 libmseed functions. The build also exports `malloc`
and `free`, for a total of 90 functions.

| Group                      | Functions                                                                             |
| -------------------------- | ------------------------------------------------------------------------------------- |
| Record parsing and packing | `msr3_*`, `ms_decode_data`, `ms_parse_raw2`, `ms_parse_raw3`                          |
| Selections                 | `ms3_*select*`, `ms3_readselectionsfile`, `ms3_freeselections`, `ms3_printselections` |
| Trace lists                | `mstl3_*`                                                                             |
| File I/O                   | `ms3_read*`, `msr3_writemseed`, `mstl3_writemseed`                                    |
| SID and channel conversion | `ms_sid2nslc_n`, `ms_nslc2sid`, `ms_seedchan2xchan`, `ms_xchan2seedchan`              |
| Extra headers              | `mseh_*`                                                                              |
| Time and leap seconds      | `ms_*time*`, `ms_doy2md`, `ms_md2doy`, `ms_readleapseconds`, `ms_readleapsecondfile`  |
| Encoding, errors and CRC   | `ms_samplesize`, `ms_encoding*`, `ms_errorstr`, `ms_crc32c`                           |
| Memory                     | `malloc`, `free`                                                                      |

## Not Exported

| Symbols                                                                          |
| -------------------------------------------------------------------------------- |
| `msr3_pack`, `mstl3_pack`, `mstl3_pack_segment`, `mstl3_pack_ppupdate_flushidle` |
| `ms_nstime2timestr`, `ms_nstime2timestrz`, `ms_sid2nslc`                         |
| `ms3_url_*`, `libmseed_url_support`                                              |
| `ms_rlog*`                                                                       |
| `ms_strncpclean`, `ms_strncpcleantail`, `ms_strncpopen`                          |
| `ms_bigendianhost`, `lmp_systemtime`, `leapsecondlist`                           |
| `libmseed_memory_prealloc`, `libmseed_memory`, `libmseed_prealloc_block_size`    |
