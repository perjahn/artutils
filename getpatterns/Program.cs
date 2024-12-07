using System;
using System.IO;
using System.Linq;
using System.Text.Json;
using System.Text.Json.Nodes;

class Program
{
    static int Main(string[] args)
    {
        if (args.Length == 0)
        {
            Console.WriteLine("Usage: getpatterns <filename>");
            return 1;
        }

        var filename = args[0];

        if (!File.Exists(filename))
        {
            Console.WriteLine("File not found: '{filename}'");
            return 1;
        }

        return GetPatterns(filename) ? 0 : 1;
    }

    static bool GetPatterns(string filename)
    {
        var json = File.ReadAllText(filename);

        var permissions = JsonSerializer.Deserialize<JsonObject[]>(json);
        if (permissions == null)
        {
            Console.WriteLine("No target patterns found in file.");
            return false;
        }

        foreach (var permission in permissions)
        {
            permission.TryGetPropertyValue("name", out var namenode);
            var permissionname = namenode?.ToString() ?? string.Empty;

            if (permission.TryGetPropertyValue("resources", out var node) && node is JsonObject resources &&
                resources.TryGetPropertyValue("artifact", out node) && node is JsonObject artifact &&
                artifact.TryGetPropertyValue("targets", out node) && node is JsonObject targets)
            {
                foreach (var target in targets)
                {
                    if (target.Value is JsonObject patterns)
                    {
                        string? includePatternsString = null;
                        var includeCount = 0;
                        if (patterns.TryGetPropertyValue("include_patterns", out node) && node is JsonArray includePatterns)
                        {
                            includeCount = includePatterns.Count;
                            includePatternsString = string.Join(", ", includePatterns.Select(p => $"'{p?.ToString() ?? string.Empty}'"));
                        }
                        else
                        {
                            Console.WriteLine($"Missing include_patterns in permission target: {target.Key}");
                        }

                        string? excludePatternsString = null;
                        var excludeCount = 0;
                        if (patterns.TryGetPropertyValue("exclude_patterns", out node) && node is JsonArray excludePatterns)
                        {
                            excludeCount = excludePatterns.Count;
                            excludePatternsString = string.Join(", ", excludePatterns.Select(p => $"'{p?.ToString() ?? string.Empty}'"));
                        }
                        else
                        {
                            Console.WriteLine($"Missing exclude_patterns in permission target: '{permissionname}'");
                        }

                        var s1 = string.Empty;
                        if (includePatternsString == null || includePatternsString != "'**'")
                        {
                            s1 = $" IncludePattern ({includeCount}): {includePatternsString}";
                        }

                        var s2 = string.Empty;
                        if (excludePatternsString == null || excludePatternsString != string.Empty)
                        {
                            s2 = $" ExcludePattern ({excludeCount}): {excludePatternsString}";
                        }

                        if (includePatternsString == null || includePatternsString != "'**'" || excludePatternsString == null || excludePatternsString != string.Empty)
                        {
                            var separator = s1 != string.Empty && s2 != string.Empty ? "," : string.Empty;
                            Console.WriteLine($"'{permissionname}' / '{target.Key}':{s1}{separator}{s2}");
                        }
                    }
                }
            }
        }

        return true;
    }
}
