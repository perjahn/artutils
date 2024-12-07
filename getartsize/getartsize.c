#include <stdio.h>
#include <stdlib.h>
#include <dirent.h>
#include <sys/stat.h>
#include <string.h>
#include <errno.h>

char path[4096];

struct folderinfo
{
  long size;
  int files;
};

struct folderinfo du()
{
  struct folderinfo fi;
  fi.size = 0;
  fi.files = 0;

  DIR *dir = opendir(path);
  if (dir == NULL)
  {
    printf("Couldn't open dir: '%s'\n", path);
    return fi;
  }

  char *end = path + strlen(path);

  struct dirent *entry;
  while ((entry = readdir(dir)) != NULL)
  {
    if (!strcmp(entry->d_name, ".") || !strcmp(entry->d_name, ".."))
    {
      continue;
    }

    *end = '/';
    strcpy(end + 1, entry->d_name);

    if (entry->d_type & DT_DIR)
    {
      if ((entry->d_name[0] != 0 && entry->d_name[1] != 0 && entry->d_name[2] == 0) || !strcmp(entry->d_name, "_pre"))
      {
        struct folderinfo subfi = du();

        double kb = subfi.size / 1024.0;
        double mb = kb / 1024;
        double gb = mb / 1024;

        printf("'%s': %ld (%.1fkb, %.1fmb, %.1fgb). Files: %d\n", path, subfi.size, kb, mb, gb, subfi.files);
        fi.size += subfi.size;
        fi.files += subfi.files;
      }
      else
      {
        printf("'%s': Not recursing into folder.\n", path);
      }
    }
    else
    {
      *end = '/';
      strcpy(end + 1, entry->d_name);

      struct stat st;
      if (stat(path, &st) == -1)
      {
        printf("Error: stat '%s'\n", path);
      }
      else
      {
        fi.size += st.st_size;
      }
      fi.files++;
    }

    *end = 0;
  }

  closedir(dir);

  return fi;
}

int main(int argc, char *argv[])
{
  if (argc < 1 || argc > 2)
  {
    printf("Usage: getartsize [folder]\n");
    return 1;
  }

  strcpy(path, argc < 2 ? "." : argv[1]);

  struct folderinfo fi = du();

  double kb = fi.size / 1024.0;
  double mb = kb / 1024;
  double gb = mb / 1024;
  double tb = gb / 1024;

  printf("Total: %ld (%.1fkb, %.1fmb, %.1fgb, %.1ftb). Files: %d\n", fi.size, kb, mb, gb, tb, fi.files);

  return 0;
}
