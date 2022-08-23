break # 1

def my_function()
    break # 2
end

gen my_generator()
    break # 3
end

for a in range(100)
    break
end

for value in range(2000)
    gen __anonymous()
        break # 4
    end
end

if false
    break # 5
end