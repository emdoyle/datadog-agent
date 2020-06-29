#ifndef _DEFS_H_
#define _DEFS_H_

#include "../../ebpf/c/bpf_helpers.h"

#if USE_SYSCALL_WRAPPER == 1
  #define SYSCALL_PREFIX "__x64_sys_"
  #define SYSCALL_KPROBE(syscall) SEC("kprobe/" SYSCALL_PREFIX #syscall) int kprobe__sys_##syscall(struct pt_regs *ctx)
  #define SYSCALL_KRETPROBE(syscall) SEC("kretprobe/" SYSCALL_PREFIX #syscall) int kretprobe__sys_##syscall(struct pt_regs *ctx)
#else
  #define SYSCALL_PREFIX "SyS_"
  #define SYSCALL_KPROBE(syscall) SEC("kprobe/" SYSCALL_PREFIX #syscall) int kprobe__sys_##syscall(struct pt_regs *ctx)
  #define SYSCALL_KRETPROBE(syscall) SEC("kretprobe/" SYSCALL_PREFIX #syscall) int kretprobe__sys_##syscall(struct pt_regs *ctx)
#endif

#define TTY_NAME_LEN 64

# define printk(fmt, ...)						\
		({							\
			char ____fmt[] = fmt;				\
			bpf_trace_printk(____fmt, sizeof(____fmt),	\
				     ##__VA_ARGS__);			\
		})

enum event_type
{
    EVENT_MAY_OPEN = 1,
    EVENT_VFS_MKDIR,
    EVENT_VFS_LINK,
    EVENT_VFS_RENAME,
    EVENT_VFS_UNLINK,
    EVENT_VFS_RMDIR,
    EVENT_VFS_CHMOD,
    EVENT_VFS_CHOWN,
    EVENT_VFS_UTIME,
    EVENT_EXEC,
};

struct event_t {
    u64 type;
    u64 timestamp;
    u64 retval;
};

struct process_data_t {
    // Process data
    u64  pidns;
    char comm[TASK_COMM_LEN];
    char tty_name[TTY_NAME_LEN];
    u32  pid;
    u32  tid;
    u32  uid;
    u32  gid;
};

struct bpf_map_def SEC("maps/events") events = {
    .type = BPF_MAP_TYPE_PERF_EVENT_ARRAY,
    .key_size = sizeof(__u32),
    .value_size = sizeof(__u32),
    .max_entries = 1024,
    .pinning = 0,
    .namespace = "",
};

#define send_event(ctx, event) \
    bpf_perf_event_output(ctx, &events, bpf_get_smp_processor_id(), &event, sizeof(event))

#endif